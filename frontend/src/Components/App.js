import React, { useState, useEffect } from "react";
import ReactDOM from "react-dom";
import styled from "styled-components";
import { NavBar } from "./NavBar.js";
import { useNavDropMenu } from "./Hooks/useNavDropMenu.js";
import { NavDropMenu } from "./NavDropMenu.js";

import { useModal } from "./Hooks/useModal.js";
import { Modal } from "./Modal.js";
import { useWs } from "./Hooks/useWS.js";

const StyledContentArea = styled.div`
	position: relative;
	z-index: 1;
	
	width: 120rem;
	height: 100vh;
	margin: 0 auto;
	
	display: grid;
	grid-template-columns: 32rem 86.5rem;
	grid-auto-rows: min-content;
	grid-gap: 1.5rem;
	
	padding-top: 8rem;
	
	& #content-area-main-display-content {
		position: relative;
	}	
		
`;


function App() {
  const ws = useWs();
  const ndm = useNavDropMenu();
  const modal = useModal();

  useEffect(() => {
		if(ws.rs === 1) {
			let storedJwt = window.localStorage.getItem('Pr0conJwt');
			if(storedJwt !== null) {
				let psjwt = JSON.parse(atob(storedJwt.split('.')[1]));
				let exp = new Date(psjwt['exp'] * 1000).toUTCString();
				let now = new Date(Date.now()).toUTCString();
				console.log(now);
				console.log(exp);
				if(exp > now) {
					console.log('Stored Jwt Good');
					ws.request(storedJwt,'validate-stored-jwt-token','noop');
				}
				if(exp < now) {
					ws.setLoading(false);
					window.localStorage.removeItem('Pr0conJwt'); 
				}
			} else if (storedJwt === null) {
				ws.setLoading(false);
			}
		}
	},[ws.rs]);


  const doLogOut = async() => {
		ws.setJwt('^vAr^');
		ws.setUser(null);
		ws.setVerifiedJwt(null)
		ws.setValidCredentials(null);
		window.localStorage.removeItem('Pr0conJwt'); 
  };

//   useEffect(() => {
// 	//used to trigger contenat area render
// 	console.log('user object action taken');
// },[ws.user]);

useEffect(() => {
	modal.setModalShowing(false);
},[ws.toggleModal]);


  useEffect(() => {
	ws.setValidCredentials(null);
},[modal.modalShowing])	


  
  return (
    <div>
      <NavBar {...ndm} {...modal} loading= {ws.loading} validjwt = {ws.verifiedJwt} />
      <NavDropMenu {...ndm} doLogOut={doLogOut} />


	  {/* {  ws.loading === false &&
				<StyledContentArea onMouseEnter={(e) => {ndm.setNavDropMenuPosX(-320)}}>
					test
				</StyledContentArea>
			} */}
			
			
	


      {modal.modalShowing && ( <Modal {...modal} validjwt={ws.verifiedJwt} validcreds={ws.validCredentials} request={ws.request} userAvail={ws.userAvail} setUserAvail={ws.setUserAvail} /> )}
    </div>
  );
}

if (document.getElementById("react_root")) {
  ReactDOM.render(<App />, document.getElementById("react_root"));
}
