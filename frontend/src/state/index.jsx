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
        if (book.id === action.payload.book.id) {
          return action.payload.book;
        }
        return book;
      });
      state.books = updatedBooks;
    },
    incBoughtBook: (state) => {
      state.bought = state.bought + 1;
    },
  },
});

export const { setLogin, setLogout, setBooks, setBook, incBoughtBook } =
  authSlice.actions;
export default authSlice.reducer;
