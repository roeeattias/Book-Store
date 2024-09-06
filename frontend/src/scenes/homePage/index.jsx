import Author from "components/Author";
import AuthorProfile from "components/AuthorProfile";
import Books from "components/Books";
import Search from "components/Search";
import User from "components/User";
import { useState } from "react";
import { useSelector } from "react-redux";

const HomePage = () => {
  const user = useSelector((state) => state.user);
  const token = useSelector((state) => state.token);
  const [visitedAuthor, setVisitedAuthor] = useState(null);

  return (
    <div className="grid grid-cols-3 w-full">
      {visitedAuthor === null ? (
        <></>
      ) : (
        <AuthorProfile author={visitedAuthor} />
      )}
      <div
        className={`flex m-4 w-full justify-center ${
          visitedAuthor === null ? "" : "filter blur-sm"
        }`}
      >
        {user === null ? (
          <User />
        ) : (
          <Author setVisitedAuthor={setVisitedAuthor} />
        )}
      </div>
      <div
        className={`flex m-4 justify-center p-4 max-h-screen overflow-y-auto ${
          visitedAuthor === null ? "" : "filter blur-sm"
        }`}
      >
        <Books />
      </div>
      <div
        className={`flex m-4 justify-center p-4 max-h-screen overflow-y-auto${
          visitedAuthor === null ? "" : "filter blur-sm"
        }`}
      >
        <Search setVisitedAuthor={setVisitedAuthor} />
      </div>
    </div>
  );
};

export default HomePage;
