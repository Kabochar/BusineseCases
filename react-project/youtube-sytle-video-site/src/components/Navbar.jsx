import { Stack } from "@mui/material";
import React from 'react';
import { Link } from "react-router-dom";
import Logo from "../asserts/favicon.ico";
import SearchBar from './SearchBar';


const Navbar = () => (
  <Stack direction="row" alignContent="center" p={2} sx={{ position: "sticky", background: '#000', top: 0, justifyContent: "space-between" }}>
    <Link to="/" style={{ display: "flex", alignItems: "center" }}>
      <img src={Logo} alt="logo" height={45} />
    </Link>
    <SearchBar />
  </Stack>
);

export default Navbar;