import React, { useState, useEffect } from "react";
import styled from "styled-components";
import axios from "axios";

export function useNavDropMenu() {
  const [navDropMenuPosX, setNavDropMenuPosX] = useState(-320);
  const [navDropMenuType, setNavDropMenuType] = useState("");

  return {
    navDropMenuPosX,
    setNavDropMenuPosX,
    navDropMenuType,
    setNavDropMenuType,
  };
}
