import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  user: null,
  token: null,
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
      const updatedBooks = state.Books.map((book) => {
        if (book._id === action.payload.book._id) return action.payload.Book;
        return book;
      });
      state.Books = updatedBooks;
    },
  },
});

export const { setLogin, setLogout, setBooks, setBook } = authSlice.actions;
export default authSlice.reducer;
