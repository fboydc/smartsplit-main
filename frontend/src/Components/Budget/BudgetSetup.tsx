import React, { useContext, useEffect ,useCallback, useState, ChangeEvent,KeyboardEvent } from "react";
import Context from "../../Context";
import ExpandableTable from "../Table/ExpandableTable";
import styles from "./budget.module.scss";
import { set } from "immer/dist/internal";
import { toast, ToastContainer, Bounce } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import StaticTable from "../Table/StaticTable";
import { Expense, Income, Allocation, AllocationGroup } from "../../models/types";
import { v4 as uuidv4} from 'uuid';





const Columns = [
  {
    Header: "Name",
    type: "text",
  },
  {
    Header: "Category",
    type: "select",
  },
  {
    Header: "Amount",
    accessor: "number",
  },
  {
    Header: "Allocation",
    accessor: "percent",
  },
]


const BudgetSetup = () => {


  const [categories, setCategories] = useState([]);
  const [monthlyIncome, setMonthlyIncome] = useState({value: ""});
  const [payFrequencies, setPayFrequencies] = useState([{id: 1, name: "Weekly"}, {id: 2, name: "Bi-Weekly"}, {id: 3, name: "Monthly"}]);
  const [payFrequency, setPayFrequency] = useState(2);
  const [allocatedAmt, setAllocatedAmt] = useState(0);
  const [needs, setNeeds] = useState([{ id: 1, name: "", amount: "", category: ""}]);
  const [wants, setWants] = useState([{ id: 1, name: "", amount: "", category: ""}]);
  const [debts, setDebts] = useState([{ id: 1, name: "", amount: "", category: ""}]);
  const [allocGroup, setAllocGroup] = useState<AllocationGroup[]>([]);
  const [totalNeedsPct, setTotalNeedsPct] = useState(0);
  const [totalNeedsAmt, setTotalNeedsAmt] = useState("");
  const [totalWantsPct, setTotalWantsPct] = useState(0);
  const [totalWantsAmt, setTotalWantsAmt] = useState("");
  const [allocatedPct, setAllocatedPct] = useState(0);
  const [totalAllocated, setTotalAllocated] = useState(0);
  const [totalDebtsPct, setTotalDebtsPct] = useState(0);
  const [totalDebtsAmt, setTotalDebtsAmt] = useState("");
  const [totalSavingsPct, setTotalSavingsPct] = useState(0);
  const [totalSavingsAmt, setTotalSavingsAmt] = useState("");



  

  const { dispatch, sessionToken, user_id, totalIncome } =
  useContext(Context);


  const getCategories = useCallback(async () => {
    const response = await fetch("/api/categories", {method: "GET",headers: {
        "Content-Type": "application/json",
        "Authorization": sessionToken,
      }
    })

    if (!response.ok) {
        return undefined;
    }

    const data = await response.json();

    return data.categories;

   }, [dispatch])

   const getDashboardInfo = useCallback(async () => {
    const response = await fetch("/api/budget?user_id="+user_id, {method: "GET",headers: {
        "Content-Type": "application/json",
        "Authorization": sessionToken,
      }
    })

  

    if (!response.ok) {
        return undefined;
    }

    const data = await response.json();
    console.log("RESPONSE - Budget", JSON.stringify(data))
    

    return data.budget;


  }, [dispatch])


    
   const formatNumber = (amount: string) => {
       //return number.replace
       return amount.replace(/\D/g, "").replace(/\B(?=(\d{3})+(?!\d))/g, ",")
   }

   const formatCurrency = (amount: string):string => {

    if (amount === "") {
      return "";
    }

    if (amount.indexOf(".") > 0) {
       

        var decimal_pos = amount.indexOf(".");

        var left_side = amount.substring(0, decimal_pos);
        var right_side = amount.substring(decimal_pos);

        left_side = formatNumber(left_side);
        right_side = formatNumber(right_side).substring(0, 2);

        amount = "$" + left_side + "." + right_side


      } else {

        amount = formatNumber(amount);
        amount = "$" + amount;
      }

     //setMonthlyIncome({value: amount});

     return amount
   }

   const transformData = (data: any) => {
        data.sort((a: any , b: any)=> a.name.localeCompare(b.name));
      return data
   }


   const reconvertToCurrency = (amount: string = "$0.00"): number => {

    if (typeof amount !== 'string') {
      console.error(`Expected a string but received ${typeof amount}`);
      return 0.00;
    }

      var cleanedString = amount.replace(/[$,\s]/g, '');
      // Return 0 if the cleaned string is empty
      if (cleanedString === "") {
        return 0.00;
      }
      // Use parseFloat to convert the cleaned string to a number
      const number = parseFloat(cleanedString);
      return number;
   }


   const updateMonthlyIncome = (incomeAmt: string) => { 
        allocGroup.forEach((group) => {
          group.allocation_pct = (group.allocation_total / reconvertToCurrency(incomeAmt)) * 100;
        })

        setAllocGroup(allocGroup);
        updateTotalAllocation(reconvertToCurrency(incomeAmt), allocGroup);

        dispatch({ type: "SET_STATE", state: { totalIncome: reconvertToCurrency(incomeAmt) }});

   }


      const calculatePercentages = (totalNeeds: number, totalWants: number, totalDebts: number, income: number, totalAllocated: number) => {
        //const reconvertedValue = reconvertToCurrency(monthlyIncome.value);

        if(income === 0) {
          setAllocatedPct(0);
          setTotalNeedsPct(0);
          setTotalWantsPct(0) ;
          setTotalDebtsPct(0);
          setTotalSavingsPct(0);
          return
        }

       // const totalAllocated = needs.reduce((total, need) => total + reconvertToCurrency(need.amount), 0) + debts.reduce((total, debt) => total + reconvertToCurrency(debt.amount), 0) + wants.reduce((total, want) => total + reconvertToCurrency(want.amount), 0);

        const allocatedPct = (totalAllocated/income) * 100;
        setAllocatedPct(parseFloat(allocatedPct.toFixed(2)));

        const totalNeedsPct = (totalNeeds / income) * 100;
        setTotalNeedsPct(parseFloat(totalNeedsPct.toFixed(2)));

        const totalWantsPct = (totalWants / income) * 100;
        setTotalWantsPct(parseFloat(totalWantsPct.toFixed(2)));
      
        const totalDebtsPct = (totalDebts / income) * 100;
        setTotalDebtsPct(parseFloat(totalDebtsPct.toFixed(2)));

        const totalSavingsPct = ((income - totalAllocated) / income) * 100;
        setTotalSavingsPct(parseFloat(totalSavingsPct.toFixed(2)));

      }

      const calculateAmts =  (totalNeeds: number, totalWants: number, totalDebts: number, income: number, totalAllocated: number) => {
          
        
        setTotalNeedsAmt(formatCurrency(totalNeeds.toString()));
        setTotalWantsAmt(formatCurrency(totalWants.toString()));
        setTotalDebtsAmt(formatCurrency(totalDebts.toString()));

        var remainingAmt = income - totalAllocated;
        console.log("income in calc amts", income)
        console.log("remaining amt", remainingAmt)
        setTotalSavingsAmt(formatCurrency(remainingAmt.toString()));

      }
      /*
      useEffect(() => {

        var totalNeeds = needs.reduce((total, need) => total + reconvertToCurrency(need.amount), 0);
        var totalWants = wants.reduce((total, want) => total + reconvertToCurrency(want.amount), 0);
        var totalDebts = debts.reduce((total, debt) => total + reconvertToCurrency(debt.amount), 0);
        var totalIncome = reconvertToCurrency(monthlyIncome.value);
        var totalAllocated = needs.reduce((total, need) => total + reconvertToCurrency(need.amount), 0) + debts.reduce((total, debt) => total + reconvertToCurrency(debt.amount), 0) + wants.reduce((total, want) => total + reconvertToCurrency(want.amount), 0);
        calculatePercentages(totalNeeds, totalWants, totalDebts, totalIncome, totalAllocated);
        calculateAmts(totalNeeds, totalWants, totalDebts, totalIncome, totalAllocated)
      }, [needs, wants, debts, monthlyIncome]);*/

   

   const handleSave = () => {
  

    


    toast.success('🤘 Budget Saved!', {
      position: "top-right",
      autoClose: 5000,
      hideProgressBar: false,
      closeOnClick: false,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      theme: "light",
      transition: Bounce,
      });
   }


   const handleUpdateExpenses = (groupIndex: number, updatedExpenses: Expense[]) => { 

      console.log("Expenses", updatedExpenses);
   
      console.log("allocation_pct", (updatedExpenses.reduce((total, expense) => Number(total) + Number(expense.amount), 0) / totalIncome) * 100)
      const updatedAllocGroup = [...allocGroup];
      updatedAllocGroup[groupIndex] = {
        ...updatedAllocGroup[groupIndex],
        expenses: updatedExpenses,
        allocation_total: updatedExpenses.reduce((total, expense) => total + Number(expense.amount), 0),
        allocation_pct: (updatedExpenses.reduce((total, expense) => total + Number(expense.amount), 0) / totalIncome) * 100,
      }
       setAllocGroup(updatedAllocGroup);
       updateTotalAllocation(totalIncome, updatedAllocGroup)
      //setAllocatedPct(updateTotalAllocation(totalIncome, updatedAllocGroup));

   }


   const updateTotalAllocation = (total_income: number, allocation_groups: AllocationGroup[]) => {

      var total = 0.0;
      var pct = 0.0;
      //console.log("alloc groups", allocGroup);
      allocation_groups.forEach((group) => {
        total =  total + group.allocation_total;
      })
      pct = (total / total_income) * 100;

      setAllocatedPct(pct)
      setTotalAllocated(total);

      //return pct;
   }


  const createExpenses = (responseExpenseData: any): Expense[] => {
     return responseExpenseData.map((expense: any) => {
        return {
          id: expense.id,
          description: expense.description,
          amount: expense.amount,
          category: expense.category,
          allocation_type: expense.allocation_type,
        }
      })
  }
   

  const createIncomes = (responseIncomeData: any): Income[] => {
    return responseIncomeData.map((income: any) => {
      return {
        id: income.id,
        amount: income.amount,
        frequency: income.frequency,
      };
    });
  };

  const createAllocations = (responseAllocationData: any): Allocation[] => {
    return responseAllocationData.map((allocation: any) => {
        return {
          type: allocation.allocation_type,
          description: allocation.allocation_description,
          factor: allocation.allocation_factor,
        }
      })
  } 

  const createAllocationGroups = (allocations: Allocation[], expenses: Expense[], income_total: number): AllocationGroup[] => {
    const allocationGroups: AllocationGroup[] = [];
    var allocatedAmt = 0;

    allocations.forEach((allocation) => {
      const group: AllocationGroup = {
        allocation_type: allocation.description,
        allocation_total: 0,
        allocation_pct: 0,
        current_allocation: 0,
        expenses: [],
      };

      expenses.forEach((expense) => {
        if (expense.allocation_type === allocation.type) {
          group.allocation_total += expense.amount;
          group.expenses.push(expense);
        }
      });

      group.allocation_pct = (group.allocation_total / income_total) * 100;
      allocatedAmt += group.allocation_total;

      allocationGroups.push(group);
    });

    /*const savingAllocations: AllocationGroup = {
      allocation_type: "Savings",
      allocation_total: income_total - allocatedAmt,
      allocation_pct: ((income_total - allocatedAmt)/ income_total) * 100,
      current_allocation: 0,
      expenses: [],
    }*/

  /*  const savingsBucket: Expense = {
      id: expenses[expenses.length-1].id + 1,
      description: "Savings Bucket",
      amount: income_total - allocatedAmt,
      category: "Expense Bucket",
      allocation_type: "savings-" + uuidv4(),
    };

    savingAllocations.expenses.push(savingsBucket);
    allocationGroups.push(savingAllocations);*/

    return allocationGroups;
  }

  const truncateDecimals = (number: number) => {
    console.log("number before truncating:" + number)
    const decimalPlaces = 2;
    const factor = Math.pow(10, decimalPlaces);
    return Math.round(number * factor) / factor;
  }





    useEffect(()=> {
      const init = async () => {
        const categories = await getCategories();
        const budget = await getDashboardInfo();

        if (categories === undefined) {
            toast.error('Error Fetching Categories!', 
            {
              position: "top-right",
              autoClose: 5000,
              hideProgressBar: false,
              closeOnClick: false,
              pauseOnHover: true,
              draggable: true,
              progress: undefined,
              theme: "dark",
              transition: Bounce,
            }
          );
        } else {
          setCategories(transformData(categories));
          setPayFrequency(2);
          setNeeds([{ id: 1, name: "", amount: "", category: categories[0].ID}]);
          setWants([{ id: 1, name: "", amount: "", category: categories[0].ID}]);
          setDebts([{ id: 1, name: "", amount: "", category: categories[0].ID}]); 
        }
      
        if (budget === undefined) {
            toast.error('Error Fetching Budget!', 
            {
              position: "top-right",
              autoClose: 5000,
              hideProgressBar: false,
              closeOnClick: false,
              pauseOnHover: true,
              draggable: true,
              progress: undefined,
              theme: "dark",
              transition: Bounce,
            }
          );
        } else {
          const incomes = createIncomes(budget.incomes);
          const expenses = createExpenses(budget.expenses);
          const allocations = createAllocations(budget.allocations);

          const income_total = incomes.reduce((total, income) => total + income.amount, 0);
          console.log("Total Income", income_total);
          const allocationGroups = createAllocationGroups(allocations, expenses, income_total);
          console.log("Allocation Groups", allocationGroups);

          setAllocGroup(allocationGroups);
          updateTotalAllocation(income_total, allocationGroups);

          //Should call filter expenses by allocations here
          
          setPayFrequency(budget.pay_frequency);
          dispatch({ type: "SET_STATE", state: { incomes: incomes, totalIncome: income_total, expenses: expenses, allocations: allocations }});
        }
        
      } 

      init();
    }, [])

    return (
      <div className={`${styles['col-md-10']} ${styles['col-lg-10']} ${styles['container']}`}>
        <h2>Budgeting Strategy</h2>
        <div>
          <h3>Income Distribution</h3>
          <ToastContainer position="top-right" autoClose={5000} hideProgressBar={false} newestOnTop={false} closeOnClick={false} rtl={false} pauseOnFocusLoss draggable pauseOnHover theme="light" transition={Bounce} />
          <hr />
          <p>Enter your monthly income and we will help you allocate it to your needs, wants, and debt repayment.</p>
          <br />
          <div className={styles.row}>
            <div className={`${styles['col-md-2']} ${styles['col-lg-2']}  ${styles['form-input']}`}>
              <label>Monthly Income</label>
              <input type="text" placeholder="$5,000" onChange={(e)=>{updateMonthlyIncome(e.target.value)}} value={formatCurrency(totalIncome + "")} className={styles.inputIncome}/>
            </div>
            <div className={`${styles['col-md-2']} ${styles['col-lg-2']}  ${styles['form-input']}`}>
              <label>Pay Frequency</label>
              <select className={styles.inputIncome} value={payFrequency} onChange={(e) => setPayFrequency(parseInt(e.target.value))}>
                {payFrequencies.map((frequency) => (
                  <option key={frequency.id} value={frequency.id}>{frequency.name}</option>
                ))}
              </select>
            </div>
             <div className={`${styles['col-lg-2']}`}>
                <label>Total Non-Savings Allocation</label>
                <p>{truncateDecimals(allocatedPct)}%</p>
            </div>
            <div className={`${styles['col-lg-2']}`}>
              <label>Unallocated Amt</label>
              <p>{totalIncome - totalAllocated}</p>
            </div>
          </div>
        </div>
        <hr />
        <div>
          <h3>Income Allocation</h3>
          <hr />
              {
                allocGroup.map((group, index) => (
                  <div key={index}>
                    <h4>{group.allocation_type} <span>{truncateDecimals(group.allocation_pct)}%</span></h4>
                    <ExpandableTable 
                      categories={categories} 
                      fields={group.expenses} 
                      setFields={(updatedExpenses)=> handleUpdateExpenses(index, updatedExpenses)} 
                      formatCurrency={formatCurrency}
                      subtotal={formatCurrency(group.allocation_total + "")}
                      />
                  </div>
                ))
              }
          <hr />
          <div>
            <button className={styles.tableButton} onClick={handleSave}>Save</button>
          </div>
          <br/>
        </div>
      </div>
    );
};


export default BudgetSetup;
