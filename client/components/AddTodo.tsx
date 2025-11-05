import { useState } from "react";
import { useForm } from "@mantine/form"; 
import { Button, Group, Modal, Textarea, TextInput } from "@mantine/core";
import { ENPOINT } from "../src/config";
import type {ToDo} from "../types/Todo"
import type { KeyedMutator } from "swr";



const AddTodo = ({mutate}:{ mutate:KeyedMutator<ToDo[]>}) => {
  const [open, setOpen] = useState(false);

  const form = useForm({
    initialValues: {
      title: "",
      body: "",
    },
  });

  const createTodo = async (values: { title: string; body: string }) =>{
    const updated = await fetch(`${ENPOINT}/api/v1/todos`,{
        method: 'POST',
        headers:{
            "Content-Type": "application/json"
        },
        body: JSON.stringify(values)
    }).then((r) => r.json())

    mutate(updated)
    form.reset();
    setOpen(false);

  }

  return (
    <>
      <Modal 
        opened={open} 
        onClose={() => setOpen(false)} 
        title="Create todo"
        centered
      >
       <form onSubmit={form.onSubmit(createTodo)}>
        <TextInput
        required
        mb={12}
        label="ToDo"
        placeholder="What you going to do?"
        {...form.getInputProps("title")}
         />
        <Textarea
         required
        mb={12}
        label="Body"
        placeholder="Tell me more..."
        {...form.getInputProps("body")}
        />
        <Button type="submit" mt="md">Create todo</Button>
       </form>
      </Modal>

      <Group justify="center">
        <Button fullWidth mt="md" onClick={() => setOpen(true)}>
          ADD TODO
        </Button>
      </Group>
    </>
  );
};

export default AddTodo;
