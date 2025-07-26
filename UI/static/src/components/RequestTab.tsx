import {
  Button,
  Container,
  Heading,
  Input,
  NativeSelect,
  Separator,
  Stack,
} from "@chakra-ui/react";

export const RequestTab = () => {
  return (
    <Container maxW={"lg"}>
      <Stack>
        <Heading textAlign={"center"}>View club requests</Heading>
        <NativeSelect.Root>
          <NativeSelect.Field placeholder="Select category">
            <option>All users</option>
            <option>Users from ITMO</option>
            <option>Users not from ITMO</option>
          </NativeSelect.Field>
          <NativeSelect.Indicator />
        </NativeSelect.Root>
        <Button>Fetch requests</Button>
        <Separator />
        <Heading textAlign={"center"}>Request approval form</Heading>
        <Input placeholder="Enter request ID" />
        <Stack direction={"row"}>
          <Button flexGrow={1} colorPalette={"red"} variant={"outline"}>
            Reject
          </Button>
          <Button flexGrow={1} colorPalette={"green"}>
            Accept
          </Button>
        </Stack>
      </Stack>
    </Container>
  );
};
