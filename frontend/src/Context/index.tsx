import { createContext, useReducer, Dispatch, ReactNode, useEffect } from "react";
import { Expense, Income, Allocation } from "../models/types";

interface QuickstartState {
  expenses: Expense[];
  allocations: Allocation[];
  linkSuccess: boolean;
  incomes: Income[];
  totalIncome: number; 
  isItemAccess: boolean;
  isPaymentInitiation: boolean;
  isUserTokenFlow: boolean;
  isCraProductsExclusively: boolean;
  linkToken: string | null;
  accessToken: string;
  userToken: string | null;
  itemId: string | null;
  isError: boolean;
  backend: boolean;
  products: string[];
  linkTokenError: {
    error_message: string;
    error_code: string;
    error_type: string;
  },
  authError: {
    error_code: string;
  }
  isAuthenticated: boolean;
  user_id: string | null;
  user: string | null;
  sessionToken: string;
  savedAllocationGroups: Allocation[];
}

const initialState: QuickstartState = {
  expenses: [],
  incomes: [],
  totalIncome: 0,
  allocations: [],
  linkSuccess: false,
  isItemAccess: true,
  isPaymentInitiation: false,
  isCraProductsExclusively: false,
  isUserTokenFlow: false,
  linkToken: "", // Don't set to null or error message will show up briefly when site loads
  userToken: null,
  accessToken: "",
  itemId: null,
  isError: false,
  backend: true,
  products: ["transactions"],
  linkTokenError: {
    error_type: "",
    error_code: "",
    error_message: "",
  },
  authError: {
    error_code: "",
  },
  isAuthenticated: false,
  user_id: null,
  user: null,
  sessionToken: "",
  savedAllocationGroups: []
};


type QuickstartAction = {
  type: "SET_STATE";
  state: Partial<QuickstartState>;
};

interface QuickstartContext extends QuickstartState {
  dispatch: Dispatch<QuickstartAction>;
}

const Context = createContext<QuickstartContext>(
  initialState as QuickstartContext
);



const { Provider } = Context;
export const QuickstartProvider: React.FC<{ children: ReactNode }> = (
  props
) => {
  const reducer = (
    state: QuickstartState,
    action: QuickstartAction
  ): QuickstartState => {
    switch (action.type) {
      case "SET_STATE":
        return { ...state, ...action.state };
      default:
        return { ...state };
    }
  };
  const [state, dispatch] = useReducer(reducer, initialState);


  const checkIfSessionExists = () => {
    // Need to implement this function
    //It should go to server and check for an active session
    //If session exists, it should return true, else false
    
  }

  useEffect(() => {
    const storedAuth = localStorage.getItem("isAuthenticated");
    const storedUser = localStorage.getItem("user");
    //use checkIfSessionExists to check if session exists
    //if it does, set isAuthenticated to true
    //else call logout
    if (storedAuth && storedUser) {
      dispatch({ type: "SET_STATE", state: { isAuthenticated: JSON.parse(storedAuth), user: storedUser } });
    }
  }, []);


  return <Provider value={{ ...state, dispatch }}>{props.children}</Provider>;
};

export default Context;
