import { createSlice, current } from "@reduxjs/toolkit";

const initialState = {
  user: null,
  token: null,
  bought: 0,
  books: [],
};

export const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setLogin: (state, action) => {
      state.user = action.payload.user;
      state.token = action.payload.token;
    },
    setLogout: (state) => {
      state.user = null;
      state.token = null;
    },
    setBooks: (state, action) => {
      state.books = action.payload.books;
    },
    setBook: (state, action) => {
      const updatedBooks = current(state.books).map((book) => {
        const bookId = book.id;
        const bookFromPayload = action.payload.book.id;
        if (bookId === bookFromPayload) {
          return action.payload.book;
        }
        return book;
      });
      state.books = updatedBooks;
    },
    addBook: (state, action) => {
      state.books = [...state.books, action.payload.book];
      state.user.publishedBooks = [
        ...state.user.publishedBooks,
        action.payload.book.id,
      ];
    },
    deleteBook: (state, action) => {
      const booksAfterDelete = state.books.filter(function (book) {
        return book.id !== action.payload.book.id;
      });
      const userPublishedBooksAfterDelete = state.user.publishedBooks.filter(
        function (book) {
          return book.id !== action.payload.book.id;
        }
      );
      state.books = booksAfterDelete;
      state.user.publishedBooks = userPublishedBooksAfterDelete;
    },
    incBoughtBook: (state) => {
      state.bought = state.bought + 1;
    },
  },
});

export const {
  setLogin,
  setLogout,
  setBooks,
  setBook,
  addBook,
  deleteBook,
  incBoughtBook,
} = authSlice.actions;
export default authSlice.reducer;
