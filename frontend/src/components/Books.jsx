import { useSelector } from "react-redux";
import { Box } from "@mui/material";
import Book from "components/Book";

const Books = () => {
  const books = useSelector((state) => state.books);
  console.log(books);
  return (
    <Box className="w-full" display="flex" flexDirection="column" gap="1.5rem">
      {books.map((book) => (
        <Book key={book.id} book={book} />
      ))}
    </Box>
  );
};

export default Books;
