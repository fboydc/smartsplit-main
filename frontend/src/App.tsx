import React, { useEffect, useContext, useCallback } from "react";

import Header from "./Components/Headers";
import Products from "./Components/ProductTypes/Products";
import Items from "./Components/ProductTypes/Items";
import Context from "./Context";
import Login from "./Components/Session/login";

import styles from "./App.module.scss";
import { CraCheckReportProduct } from "plaid";
import Main from "./Main";
import { Navigate, Routes, Route } from "react-router";
import Categories from "./Components/Budget/BudgetSetup";

const App = () => {

  const { isAuthenticated } = useContext(Context);


  return (
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={ <Main />} />
          <Route path="/budgeting" element={ <Categories />} />
        </Routes>
  );
};

export default App;
