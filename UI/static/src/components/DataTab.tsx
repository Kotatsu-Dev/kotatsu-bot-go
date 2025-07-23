import { useAPI } from "../api/api";
import { Button, Center, Container, Heading, Stack } from "@chakra-ui/react";

export const DataTab = () => {
  const api = useAPI();

  return (
    <Container>
      <Heading textAlign={"center"} pb={3}>
        Dashboard
      </Heading>
      <Center>
        <Stack w="lg">
          <Button
            onClick={() => {
              api.users.getAll();
            }}
          >
            Show club members
          </Button>
          <Button>Show newsletter subscribers</Button>
          <Button>Show events list</Button>
        </Stack>
      </Center>
    </Container>
  );
};
