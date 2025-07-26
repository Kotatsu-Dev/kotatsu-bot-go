import { Button, Container, Heading, Stack, Textarea } from "@chakra-ui/react";

export const MessagesTab = () => {
  return (
    <Container maxW="lg">
        <Stack>
            <Heading textAlign={"center"}>Message broadcasting</Heading>
            <Textarea />
            <Button>Send broadcast</Button>
        </Stack>
    </Container>
  );
};
