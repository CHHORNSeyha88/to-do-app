import { MantineProvider, Box, List, ThemeIcon } from '@mantine/core';
import useSWR from 'swr';
import AddTodo from "../components/AddTodo";
import type { ToDo } from "../types/Todo";
import { CheckCircleFillIcon } from '@primer/octicons-react';
import { ENPOINT } from './config';

// endpoint moved to ./config to avoid circular import with AddTodo

const fetcher = (url: string) =>
  fetch(`${ENPOINT}/${url}`).then((r) => r.json());

function App() {
  const { data, mutate } = useSWR<ToDo[]>('api/v1/todos', fetcher);

  const MarkasDone = async (id: number) => {
    await fetch(`${ENPOINT}/api/v1/todos/${id}/done`, { method: "PATCH" });

    // Update SWR cache locally
    mutate((currentData) => {
      if (!currentData) return [];
      return currentData.map((todo) =>
        todo.id === id ? { ...todo, done: true } : todo
      );
    }, false);
  };

  return (
    <MantineProvider>
      <Box
        style={{
          padding: "2rem",
          width: "100%",
          maxWidth: "40rem",
          margin: "0 auto",
        }}
      >
        <List spacing="xs" size="sm" mb={12} center>
          {data?.map((todo) => (
            <List.Item
              key={`todo__${todo.id}`}
              onClick={() => MarkasDone(todo.id)}
              icon={
                todo.done ? (
                  <ThemeIcon color="teal" size={24} radius="xl">
                    <CheckCircleFillIcon size={20} />
                  </ThemeIcon>
                ) : (
                  <ThemeIcon color="gray" size={24} radius="xl">
                    <CheckCircleFillIcon size={20} />
                  </ThemeIcon>
                )
              }
            >
              {todo.title}
            </List.Item>
          ))}
        </List>
        <AddTodo mutate={mutate} />
      </Box>
    </MantineProvider>
  );
}

export default App;
