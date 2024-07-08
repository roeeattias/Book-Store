import { useSelector } from "react-redux";
import { Box } from "@mui/material";
import Book from "components/Book";

const Books = ({ booksToShow }) => {
  const books = useSelector((state) => state.books);

  return (
    <Box className="w-full" display="flex" flexDirection="column" gap="1.5rem">
      {booksToShow === undefined ? (
        <>
          {books.map((book) => (
            <Book key={`home page ${book.id}`} book={book} />
          ))}
        </>
      ) : (
        <>
          {booksToShow.map((book) => (
            <Book
              key={`profile ${book.id}`}
              book={book}
              inAuthorProfile={true}
            />
          ))}
        </>
      )}
    </Box>
  );
};

export default Books;
