import { useState, type FC } from "react";
import { useDatePicker } from "@rehookify/datepicker";
import {
  Button,
  Group,
  Heading,
  IconButton,
  Input,
  NativeSelect,
  Popover,
  Portal,
  SimpleGrid,
  Stack,
  Text,
} from "@chakra-ui/react";
import { format } from "date-fns";
import { FaChevronLeft, FaChevronRight } from "react-icons/fa";

type Props = {
  onChange?: (e: Date | null) => void;
  onBlur?: () => void;
  //   value: Date | null;
  defaultValue?: Date;
  disabled?: boolean;
  name?: string;
};

export const Calendar: FC<Props> = (props) => {
  const [selectedDates, setSelectedDates] = useState<Date[]>(
    props.defaultValue ? [props.defaultValue] : []
  );
  const [offsetDate, onOffsetChange] = useState<Date>(new Date());

  const {
    data: { weekDays, calendars, time },
    propGetters: { dayButton, timeButton, addOffset, subtractOffset },
  } = useDatePicker({
    selectedDates,
    focusDate: props.defaultValue,
    onDatesChange(dates) {
      setSelectedDates(dates);
      props.onChange?.(dates.at(0) ?? null);
    },
    offsetDate,
    onOffsetChange,
    calendar: { startDay: 1 },
  });

  const { year, month, days } = calendars[0];

  return (
    <>
      <Group w="full">
        <Popover.Root
          positioning={{ placement: "bottom-start" }}
          onOpenChange={(e) => {
            if (!e.open) {
              props.onBlur?.();
            }
          }}
        >
          <Popover.Trigger asChild>
            <Input
              value={
                selectedDates.length > 0
                  ? format(selectedDates[0], "dd.MM.yyyy")
                  : "??.??.????"
              }
              readOnly
              disabled={props.disabled}
              name={props.name}
            />
          </Popover.Trigger>
          <Portal>
            <Popover.Positioner>
              <Popover.Content colorPalette={"orange"}>
                <Popover.Arrow />
                <Popover.Body>
                  <Stack direction={"row"} alignItems={"center"} mb={4}>
                    <IconButton
                      {...subtractOffset({ months: 1 })}
                      aria-label="Previous month"
                    >
                      <FaChevronLeft />
                    </IconButton>
                    <Heading flexGrow={1} textAlign={"center"}>
                      {month} {year}
                    </Heading>
                    <IconButton
                      {...addOffset({ months: 1 })}
                      aria-label="Next month"
                    >
                      <FaChevronRight />
                    </IconButton>
                  </Stack>
                  <SimpleGrid columns={7} gap={2}>
                    {weekDays.map((day) => (
                      <Text key={`${month}-${day}`} textAlign={"center"}>
                        {day}
                      </Text>
                    ))}
                    {days.map((dpDay) => (
                      <Button
                        key={dpDay.$date.toDateString()}
                        {...dayButton(dpDay)}
                        variant={dpDay.selected ? "outline" : "ghost"}
                        colorPalette={dpDay.inCurrentMonth ? "current" : "gray"}
                      >
                        {dpDay.day}
                      </Button>
                    ))}
                  </SimpleGrid>
                </Popover.Body>
              </Popover.Content>
            </Popover.Positioner>
          </Portal>
        </Popover.Root>
        <NativeSelect.Root
          onBlur={props.onBlur}
          onChange={(e) =>
            timeButton(time[(e.target as any).value]).onClick?.(e as any)
          }
        >
          <NativeSelect.Field
            defaultValue={time.findIndex((dpTime) => dpTime.selected)}
          >
            {time.map((dpTime, i) => (
              <option
                key={`${dpTime.$date.toDateString()} ${dpTime.time}`}
                value={i}
                {...timeButton(dpTime)}
              >
                {dpTime.time}
              </option>
            ))}
          </NativeSelect.Field>
          <NativeSelect.Indicator />
        </NativeSelect.Root>
      </Group>
    </>
  );
};
