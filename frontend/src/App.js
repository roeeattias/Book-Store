import "./App.css";
import { useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from "scenes/homePage";
import { useDispatch, useSelector } from "react-redux";
import { CssBaseline } from "@mui/material";
import { setBooks } from "state";
import LoginPage from "scenes/loginPage";
import SignUpPage from "scenes/signupPage";

function App() {
  const isAuth = Boolean(useSelector((state) => state.token));
  const dispatch = useDispatch();

  const fetchBooks = async () => {
    // fetching the books
    const response = await fetch("http://localhost:8080/getBooks", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    });
    const books = await response.json();
    if (books) {
      // setting the books in the state
      dispatch(
        setBooks({
          books: books,
        })
      );
    }
  };

  // Getting the book as soon as the app is rendered.
  useEffect(() => {
    fetchBooks(); // calling the functio to get the books
  }, []);

  return (
    <div className="app bg-background flex flex-row h-screen">
      <BrowserRouter>
        <CssBaseline />
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/signup" element={<SignUpPage />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
