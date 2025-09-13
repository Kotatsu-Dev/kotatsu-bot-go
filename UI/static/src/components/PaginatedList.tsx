import {
  ButtonGroup,
  Center,
  IconButton,
  Pagination,
  Stack,
} from "@chakra-ui/react";
import { useEffect, useState, type ReactNode } from "react";
import { FaChevronLeft, FaChevronRight } from "react-icons/fa";

export const PaginatedList = <T,>(props: {
  items: T[];
  pageSize: number;
  render: (e: T) => ReactNode;
}) => {
  const [page, setPage] = useState(1);

  useEffect(() => {
    setPage(1);
  }, [props.items, props.pageSize]);

  return (
    <Stack>
      <Center>
        <Pagination.Root
          count={props.items.length}
          pageSize={props.pageSize}
          page={page}
          onPageChange={(e) => setPage(e.page)}
        >
          <ButtonGroup variant="ghost" size="sm">
            <Pagination.PrevTrigger asChild>
              <IconButton>
                <FaChevronLeft />
              </IconButton>
            </Pagination.PrevTrigger>
            <Pagination.Items
              render={(page) => (
                <IconButton variant={{ base: "ghost", _selected: "outline" }}>
                  {page.value}
                </IconButton>
              )}
            />
            <Pagination.NextTrigger asChild>
              <IconButton>
                <FaChevronRight />
              </IconButton>
            </Pagination.NextTrigger>
          </ButtonGroup>
        </Pagination.Root>
      </Center>
      {props.items
        .slice(props.pageSize * (page - 1), props.pageSize * page)
        .map(props.render)}
    </Stack>
  );
};
