import { handleError, useAPI } from "../api/api";
import { Button, Container, Heading, Input, Stack } from "@chakra-ui/react";
import { useState } from "react";
import { toaster } from "./ui/toaster";

export const UsersTab = () => {
  const api = useAPI();
  const [userId, setUserId] = useState(0);

  const kickUser = async () => {
    try {
      await api.users.setMemberStatus({
        user_tg_id: userId,
        is_club_member: false,
      });
      toaster.success({
        description: "Successfully kicked user",
      });
    } catch (e) {
      handleError(e);
    }
  };

  const addUser = async () => {
    try {
      await api.users.setMemberStatus({
        user_tg_id: userId,
        is_club_member: true,
      });
      toaster.success({
        description: "Successfully added user",
      });
    } catch (e) {
      handleError(e);
    }
  };

  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>User management</Heading>
        <Input
          placeholder="Enter user's Telegram ID"
          type="number"
          onChange={(e) => setUserId(e.currentTarget.valueAsNumber)}
        />
        <Stack direction={"row"}>
          <Button
            variant={"outline"}
            colorPalette={"red"}
            flexGrow={1}
            onClick={kickUser}
          >
            Kick user from club
          </Button>
          <Button colorPalette={"green"} flexGrow={1} onClick={addUser}>
            Add user to club
          </Button>
        </Stack>
      </Stack>
    </Container>
  );
};
