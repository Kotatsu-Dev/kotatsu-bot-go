import { Button, Center, Container, Heading, Stack } from "@chakra-ui/react";

export const DataTab = () => {
  return (
    <Container>
      <Heading textAlign={'center'} pb={3}>Dashboard</Heading>
      <Center>
        <Stack w="lg">
            <Button>Show club members</Button>
            <Button>Show newsletter subscribers</Button>
            <Button>Show events list</Button>
        </Stack>
      </Center>
    </Container>
  );
};
