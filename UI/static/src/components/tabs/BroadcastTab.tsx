import { handleError, useAPI } from "../../api/api";
import {
  Badge,
  Button,
  Card,
  CloseButton,
  Combobox,
  Container,
  createListCollection,
  Dialog,
  Field,
  Flex,
  Heading,
  IconButton,
  Input,
  Listbox,
  Portal,
  Stack,
  Table,
  Text,
  Textarea,
  useFilter,
  useListCollection,
  Status,
} from "@chakra-ui/react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useDebounceValue } from "usehooks-ts";
import { toaster } from "../ui/toaster";
import { itmoStatusCollection as itmoStatus } from "../../constants/users";
import { searchUsers } from "../../lib/userSearch";
import type { User } from "../../api/users";
import type { Activity } from "../../api/activities";
import type { Roulette } from "../../api/roulettes";
import type { BroadcastResult } from "../../api/broadcast";
import { FaTimes } from "react-icons/fa";

const clubMemberStatus = createListCollection({
  items: [
    { label: "Club member", value: "club_member" },
    { label: "Not club member", value: "not_club_member" },
  ],
});

type Rule =
  | { id: string; kind: "event"; eventId: number; label: string }
  | { id: string; kind: "roulette"; rouletteId: number; label: string }
  | {
      id: string;
      kind: "filter";
      clubStatuses: string[];
      itmoStatuses: string[];
      label: string;
    };

const displayName = (user: User) =>
  user.full_name || (user.user_name ? `@${user.user_name}` : "") ||
  `Telegram ID ${user.user_tg_id}`;

export const BroadcastTab = () => {
  const api = useAPI();
  const { contains } = useFilter({ sensitivity: "base" });

  const [broadcastResults, setBroadcastResults] = useState<BroadcastResult[]>(
    [],
  );
  const [resultsOpen, setResultsOpen] = useState(false);
  const [message, setMessage] = useState("");

  const [allEvents, setAllEvents] = useState<Activity[]>([]);
  const [allRoulettes, setAllRoulettes] = useState<Roulette[]>([]);
  const [allUsers, setAllUsers] = useState<User[]>([]);

  const [selectedEventIds, setSelectedEventIds] = useState<number[]>([]);
  const [selectedRouletteIds, setSelectedRouletteIds] = useState<number[]>(
    [],
  );
  const [selectedClubStatuses, setSelectedClubStatuses] = useState<string[]>(
    [],
  );
  const [selectedItmoStatuses, setSelectedItmoStatuses] = useState<string[]>(
    [],
  );

  const [rules, setRules] = useState<Rule[]>([]);
  const [manualTgIds, setManualTgIds] = useState<Set<number>>(new Set());
  const [excludedTgIds, setExcludedTgIds] = useState<Set<number>>(new Set());

  const [userSearchInput, setUserSearchInput] = useState("");
  const [userSearch, setUserSearch] = useDebounceValue("", 500);

  const nextRuleId = useRef(0);
  const newRuleId = () => `rule-${nextRuleId.current++}`;

  const {
    collection: events,
    filter: filterEvents,
    set: setEventsCollection,
  } = useListCollection<{ label: string; value: number }>({
    initialItems: [],
    filter: contains,
  });

  const {
    collection: roulettes,
    filter: filterRoulettes,
    set: setRoulettesCollection,
  } = useListCollection<{ label: string; value: number }>({
    initialItems: [],
    filter: contains,
  });

  const loadEvents = async () => {
    const sortedEvents = (await api.activities.getAll()).sort(
      (a, b) => b.id - a.id,
    );
    setAllEvents(sortedEvents);
    setEventsCollection(
      sortedEvents.map((event) => ({ label: event.title, value: event.id })),
    );
  };

  const loadRoulettes = async () => {
    const allRoulettesData = await api.roulettes.getAll();
    setAllRoulettes(allRoulettesData);
    setRoulettesCollection(
      allRoulettesData.map((roulette) => ({
        label: roulette.theme,
        value: roulette.id,
      })),
    );
  };

  const loadUsers = async () => {
    setAllUsers(await api.users.getAll());
  };

  useEffect(() => {
    loadEvents();
    loadRoulettes();
    loadUsers();
  }, []);

  const byTgId = useMemo(
    () => new Map(allUsers.map((user) => [user.user_tg_id, user])),
    [allUsers],
  );

  const addEventRules = () => {
    const newRules: Rule[] = selectedEventIds
      .map((id) => allEvents.find((event) => event.id === id))
      .filter((event): event is Activity => !!event)
      .map((event) => ({
        id: newRuleId(),
        kind: "event",
        eventId: event.id,
        label: `Event: ${event.title}`,
      }));
    setRules((prev) => [...prev, ...newRules]);
    setSelectedEventIds([]);
  };

  const addRouletteRules = () => {
    const newRules: Rule[] = selectedRouletteIds
      .map((id) => allRoulettes.find((roulette) => roulette.id === id))
      .filter((roulette): roulette is Roulette => !!roulette)
      .map((roulette) => ({
        id: newRuleId(),
        kind: "roulette",
        rouletteId: roulette.id,
        label: `Roulette: ${roulette.theme}`,
      }));
    setRules((prev) => [...prev, ...newRules]);
    setSelectedRouletteIds([]);
  };

  const addFilterRule = () => {
    if (selectedClubStatuses.length === 0 && selectedItmoStatuses.length === 0) {
      return;
    }

    const parts: string[] = [];
    if (selectedClubStatuses.length > 0) {
      parts.push(
        selectedClubStatuses
          .map(
            (value) =>
              clubMemberStatus.items.find((item) => item.value === value)
                ?.label ?? value,
          )
          .join(", "),
      );
    }
    if (selectedItmoStatuses.length > 0) {
      parts.push(
        selectedItmoStatuses
          .map(
            (value) =>
              itmoStatus.items.find((item) => item.value === value)?.label ??
              value,
          )
          .join(", "),
      );
    }

    setRules((prev) => [
      ...prev,
      {
        id: newRuleId(),
        kind: "filter",
        clubStatuses: selectedClubStatuses,
        itmoStatuses: selectedItmoStatuses,
        label: `Filter: ${parts.join(" · ")}`,
      },
    ]);
    setSelectedClubStatuses([]);
    setSelectedItmoStatuses([]);
  };

  const removeRule = (id: string) => {
    setRules((prev) => prev.filter((rule) => rule.id !== id));
  };

  const resolveEventAudience = (eventId: number): User[] => {
    const event = allEvents.find((e) => e.id === eventId);
    if (!event) return [];
    const result: User[] = [];
    for (const participant of event.participants) {
      const user = byTgId.get(participant.user_tg_id);
      if (user) result.push(user);
    }
    return result;
  };

  const resolveRouletteAudience = (rouletteId: number): User[] => {
    const roulette = allRoulettes.find((r) => r.id === rouletteId);
    if (!roulette?.participants) return [];
    return roulette.participants.map(
      (participant) => byTgId.get(participant.user_tg_id) ?? participant,
    );
  };

  const resolveFilterAudience = (
    rule: Extract<Rule, { kind: "filter" }>,
  ): User[] =>
    allUsers.filter((user) => {
      if (rule.clubStatuses.length > 0) {
        const wantMember = rule.clubStatuses.includes("club_member");
        const wantNotMember = rule.clubStatuses.includes("not_club_member");
        if (wantMember && !wantNotMember && !user.is_club_member) return false;
        if (wantNotMember && !wantMember && user.is_club_member) return false;
      }
      if (
        rule.itmoStatuses.length > 0 &&
        !rule.itmoStatuses.includes(user.itmo_status ?? "")
      ) {
        return false;
      }
      return true;
    });

  const recipients = useMemo(() => {
    const map = new Map<number, { user: User; sources: string[] }>();

    const include = (user: User, source: string) => {
      if (excludedTgIds.has(user.user_tg_id)) return;
      const existing = map.get(user.user_tg_id);
      if (existing) {
        if (!existing.sources.includes(source)) existing.sources.push(source);
      } else {
        map.set(user.user_tg_id, { user, sources: [source] });
      }
    };

    for (const rule of rules) {
      const audience =
        rule.kind === "event"
          ? resolveEventAudience(rule.eventId)
          : rule.kind === "roulette"
            ? resolveRouletteAudience(rule.rouletteId)
            : resolveFilterAudience(rule);
      for (const user of audience) include(user, rule.label);
    }

    for (const tgId of manualTgIds) {
      const user = byTgId.get(tgId);
      if (user) include(user, "Manually added");
    }

    return Array.from(map.values());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [rules, manualTgIds, excludedTgIds, allEvents, allRoulettes, allUsers, byTgId]);

  const recipientTgIds = useMemo(
    () => new Set(recipients.map((r) => r.user.user_tg_id)),
    [recipients],
  );

  const addManualUser = (user: User) => {
    setManualTgIds((prev) => new Set(prev).add(user.user_tg_id));
    setExcludedTgIds((prev) => {
      if (!prev.has(user.user_tg_id)) return prev;
      const next = new Set(prev);
      next.delete(user.user_tg_id);
      return next;
    });
  };

  const removeRecipient = (tgId: number) => {
    setExcludedTgIds((prev) => new Set(prev).add(tgId));
    setManualTgIds((prev) => {
      if (!prev.has(tgId)) return prev;
      const next = new Set(prev);
      next.delete(tgId);
      return next;
    });
  };

  const userSearchResults = useMemo(() => {
    if (userSearch.trim().length === 0) return [];
    return searchUsers(allUsers, userSearch).slice(0, 8);
  }, [allUsers, userSearch]);

  const handleSend = async () => {
    try {
      const results = await api.broadcast.sendBroadcast({
        events: [],
        users: recipients.map(({ user }) => user.id),
        roulettes: [],
        club_member_status: null,
        itmo_status: [],
        message,
      });

      setBroadcastResults(results);
      setResultsOpen(true);

      toaster.success({
        title: "Broadcast sent successfully",
      });

      setRules([]);
      setManualTgIds(new Set());
      setExcludedTgIds(new Set());
      setMessage("");
    } catch (error) {
      handleError(error);
    }
  };

  return (
    <Container maxW={"2xl"} mb={5}>
      <Heading textAlign={"center"} pb={3}>
        Broadcast to users
      </Heading>
      <Stack gap={4}>
        <Card.Root>
          <Card.Body>
            <Stack gap={3}>
              <Text fontWeight={"medium"}>Add rule by event</Text>
              <Combobox.Root
                multiple
                closeOnSelect={false}
                collection={events}
                onInputValueChange={(e) => filterEvents(e.inputValue)}
                onValueChange={({ value }) =>
                  setSelectedEventIds(value as unknown as number[])
                }
                value={selectedEventIds as unknown as string[]}
              >
                <Combobox.Control>
                  <Combobox.Input placeholder="Type to search events" />
                  <Combobox.IndicatorGroup>
                    <Combobox.ClearTrigger />
                    <Combobox.Trigger />
                  </Combobox.IndicatorGroup>
                </Combobox.Control>
                <Portal>
                  <Combobox.Positioner>
                    <Combobox.Content>
                      <Combobox.Empty>No items found</Combobox.Empty>
                      {events.items.map((item) => (
                        <Combobox.Item item={item} key={item.value}>
                          {item.label}
                          <Combobox.ItemIndicator />
                        </Combobox.Item>
                      ))}
                    </Combobox.Content>
                  </Combobox.Positioner>
                </Portal>
              </Combobox.Root>
              <Flex justify={"space-between"} align={"center"}>
                <Text color={"fg.muted"}>
                  {selectedEventIds.length > 0
                    ? `${selectedEventIds.length} event(s) selected`
                    : "No events selected"}
                </Text>
                <Button
                  size="sm"
                  onClick={addEventRules}
                  disabled={selectedEventIds.length === 0}
                >
                  Add rule
                </Button>
              </Flex>
            </Stack>
          </Card.Body>
        </Card.Root>

        <Card.Root>
          <Card.Body>
            <Stack gap={3}>
              <Text fontWeight={"medium"}>Add rule by roulette</Text>
              <Combobox.Root
                multiple
                closeOnSelect={false}
                collection={roulettes}
                onInputValueChange={(e) => filterRoulettes(e.inputValue)}
                onValueChange={({ value }) =>
                  setSelectedRouletteIds(value as unknown as number[])
                }
                value={selectedRouletteIds as unknown as string[]}
              >
                <Combobox.Control>
                  <Combobox.Input placeholder="Type to search roulettes" />
                  <Combobox.IndicatorGroup>
                    <Combobox.ClearTrigger />
                    <Combobox.Trigger />
                  </Combobox.IndicatorGroup>
                </Combobox.Control>
                <Portal>
                  <Combobox.Positioner>
                    <Combobox.Content>
                      <Combobox.Empty>No items found</Combobox.Empty>
                      {roulettes.items.map((item) => (
                        <Combobox.Item item={item} key={item.value}>
                          {item.label}
                          <Combobox.ItemIndicator />
                        </Combobox.Item>
                      ))}
                    </Combobox.Content>
                  </Combobox.Positioner>
                </Portal>
              </Combobox.Root>
              <Flex justify={"space-between"} align={"center"}>
                <Text color={"fg.muted"}>
                  {selectedRouletteIds.length > 0
                    ? `${selectedRouletteIds.length} roulette(s) selected`
                    : "No roulettes selected"}
                </Text>
                <Button
                  size="sm"
                  onClick={addRouletteRules}
                  disabled={selectedRouletteIds.length === 0}
                >
                  Add rule
                </Button>
              </Flex>
            </Stack>
          </Card.Body>
        </Card.Root>

        <Card.Root>
          <Card.Body>
            <Stack gap={3}>
              <Text fontWeight={"medium"}>Add rule by filter</Text>
              <Stack gap={1}>
                <Text fontSize={"sm"} color={"fg.muted"}>
                  Club membership
                </Text>
                <Listbox.Root
                  collection={clubMemberStatus}
                  selectionMode="multiple"
                  value={selectedClubStatuses}
                  onValueChange={({ value }) => setSelectedClubStatuses(value)}
                >
                  <Listbox.Content>
                    {clubMemberStatus.items.map((membership) => (
                      <Listbox.Item item={membership} key={membership.value}>
                        <Listbox.ItemText>{membership.label}</Listbox.ItemText>
                        <Listbox.ItemIndicator />
                      </Listbox.Item>
                    ))}
                  </Listbox.Content>
                </Listbox.Root>
              </Stack>
              <Stack gap={1}>
                <Text fontSize={"sm"} color={"fg.muted"}>
                  ITMO status
                </Text>
                <Listbox.Root
                  collection={itmoStatus}
                  selectionMode="multiple"
                  value={selectedItmoStatuses}
                  onValueChange={({ value }) => setSelectedItmoStatuses(value)}
                >
                  <Listbox.Content>
                    {itmoStatus.items.map((status) => (
                      <Listbox.Item item={status} key={status.value}>
                        <Listbox.ItemText>{status.label}</Listbox.ItemText>
                        <Listbox.ItemIndicator />
                      </Listbox.Item>
                    ))}
                  </Listbox.Content>
                </Listbox.Root>
              </Stack>
              <Flex justify={"space-between"} align={"center"}>
                <Text color={"fg.muted"}>
                  {selectedClubStatuses.length + selectedItmoStatuses.length >
                  0
                    ? "Filter configured"
                    : "No filter configured"}
                </Text>
                <Button
                  size="sm"
                  onClick={addFilterRule}
                  disabled={
                    selectedClubStatuses.length === 0 &&
                    selectedItmoStatuses.length === 0
                  }
                >
                  Add rule
                </Button>
              </Flex>
            </Stack>
          </Card.Body>
        </Card.Root>

        {rules.length > 0 ? (
          <Card.Root>
            <Card.Body>
              <Stack gap={2}>
                <Text fontWeight={"medium"}>Active rules</Text>
                {rules.map((rule) => (
                  <Flex
                    key={rule.id}
                    justify={"space-between"}
                    align={"center"}
                    gap={2}
                  >
                    <Text>{rule.label}</Text>
                    <IconButton
                      aria-label="Remove rule"
                      size="xs"
                      variant="ghost"
                      onClick={() => removeRule(rule.id)}
                    >
                      <FaTimes />
                    </IconButton>
                  </Flex>
                ))}
              </Stack>
            </Card.Body>
          </Card.Root>
        ) : null}

        <Card.Root>
          <Card.Body>
            <Stack gap={3}>
              <Text fontWeight={"medium"}>Add individual user</Text>
              <Input
                placeholder="Search by name, ISU, phone, @username or Telegram ID"
                value={userSearchInput}
                onChange={(e) => {
                  setUserSearchInput(e.currentTarget.value);
                  setUserSearch(e.currentTarget.value);
                }}
              />
              {userSearchResults.length > 0 ? (
                <Stack gap={1}>
                  {userSearchResults.map((user) => {
                    const included = recipientTgIds.has(user.user_tg_id);
                    return (
                      <Flex
                        key={user.id}
                        justify={"space-between"}
                        align={"center"}
                        gap={2}
                        px={2}
                        py={1}
                        borderRadius={"md"}
                        cursor={"pointer"}
                        _hover={{ bg: "bg.muted" }}
                        onClick={() => addManualUser(user)}
                      >
                        <Stack gap={0}>
                          <Text>{displayName(user)}</Text>
                          <Text fontSize={"sm"} color={"fg.muted"}>
                            {[
                              user.isu,
                              user.phone_number,
                              user.user_name ? `@${user.user_name}` : "",
                            ]
                              .filter(Boolean)
                              .join(" · ")}
                          </Text>
                        </Stack>
                        {included ? (
                          <Badge colorPalette={"green"}>Added</Badge>
                        ) : null}
                      </Flex>
                    );
                  })}
                </Stack>
              ) : userSearchInput.trim().length > 0 ? (
                <Text color={"fg.muted"}>No users found</Text>
              ) : null}
            </Stack>
          </Card.Body>
        </Card.Root>

        <Card.Root>
          <Card.Body>
            <Stack gap={2}>
              <Text fontWeight={"medium"}>
                Recipients ({recipients.length})
              </Text>
              {recipients.length === 0 ? (
                <Text color={"fg.muted"}>
                  No recipients yet — add a rule or a user above
                </Text>
              ) : (
                <Table.Root size={"sm"}>
                  <Table.Header>
                    <Table.Row>
                      <Table.ColumnHeader>Full name</Table.ColumnHeader>
                      <Table.ColumnHeader>Added by</Table.ColumnHeader>
                      <Table.ColumnHeader></Table.ColumnHeader>
                    </Table.Row>
                  </Table.Header>
                  <Table.Body>
                    {recipients.map(({ user, sources }) => (
                      <Table.Row key={user.user_tg_id}>
                        <Table.Cell>{displayName(user)}</Table.Cell>
                        <Table.Cell>{sources.join(", ")}</Table.Cell>
                        <Table.Cell>
                          <IconButton
                            aria-label="Remove recipient"
                            size="xs"
                            variant="ghost"
                            onClick={() => removeRecipient(user.user_tg_id)}
                          >
                            <FaTimes />
                          </IconButton>
                        </Table.Cell>
                      </Table.Row>
                    ))}
                  </Table.Body>
                </Table.Root>
              )}
            </Stack>
          </Card.Body>
        </Card.Root>

        <Field.Root>
          <Field.Label>Message</Field.Label>
          <Textarea
            value={message}
            onChange={(e) => setMessage(e.currentTarget.value)}
          />
        </Field.Root>

        <Button
          onClick={handleSend}
          disabled={recipients.length === 0 || message.trim().length === 0}
        >
          Send broadcast to {recipients.length} user(s)
        </Button>
      </Stack>

      <Dialog.Root
        open={resultsOpen}
        onOpenChange={() => setResultsOpen(false)}
      >
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.Header>
                <Dialog.Title>Broadcast results</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <Table.Root>
                  <Table.Header>
                    <Table.Row>
                      <Table.ColumnHeader>Status</Table.ColumnHeader>
                      <Table.ColumnHeader>User</Table.ColumnHeader>
                      <Table.ColumnHeader>Error</Table.ColumnHeader>
                    </Table.Row>
                  </Table.Header>
                  <Table.Body>
                    {broadcastResults.map((result, index) => (
                      <Table.Row key={index}>
                        <Table.Cell>
                          <Status.Root
                            colorPalette={result.success ? "green" : "red"}
                          >
                            <Status.Indicator />
                          </Status.Root>
                        </Table.Cell>
                        <Table.Cell>
                          {result.user.user_name
                            ? `@${result.user.user_name}`
                            : ""}{" "}
                          ({result.user.user_tg_id})
                        </Table.Cell>
                        <Table.Cell>
                          {result.success ? "" : result.error_message}
                        </Table.Cell>
                      </Table.Row>
                    ))}
                  </Table.Body>
                </Table.Root>
              </Dialog.Body>
              <Dialog.Footer>
                <Dialog.ActionTrigger asChild>
                  <Button variant="outline">Close</Button>
                </Dialog.ActionTrigger>
              </Dialog.Footer>
              <Dialog.CloseTrigger asChild>
                <CloseButton size="sm" />
              </Dialog.CloseTrigger>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </Container>
  );
};
