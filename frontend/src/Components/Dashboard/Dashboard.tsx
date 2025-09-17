import React, { useContext } from "react";

import Endpoint from "../Endpoint";
import Context from "../../Context";
import ProductTypesContainer from "../ProductTypes/ProductTypesContainer";
import def from "ajv/dist/vocabularies/discriminator";
import BudgetSetup  from "../Budget/BudgetSetup";
import styles from "./dashboard.module.scss";
import PieChart, { PieDatum } from '../Graphs/PieCharts';

const data: PieDatum[] = [
  { id: 'a', label: 'Needs', value: 50 },
  { id: 'b', label: 'Debt & Repayment', value: 20 },
  { id: 'c', label: 'Wants', value: 10 },
  { id: 'd', label: 'Unspent', value: 20 },
];


const Dashboard = () => {

  return (
    <div>
      <div className={styles.row}>
        <div className={styles.piechartContainer}>
            <PieChart
              data={data}
              size={420}
              innerRadius={80}
              padAngleDeg={1}
              showLegend={true}
              legendPosition="right"
              onSliceClick={(d) => alert(`clicked ${d.label}`)}
            />
        </div>
        <div className={styles.progressBarContainer}>
            Allocations shown as of {new Date().toLocaleDateString()}
            <div>
                  <progress value={100} max={100}></progress> Needs
            </div>
            <div>
                  <progress value={100} max={100}></progress> Debt & Repayment
            </div>
              <div>
                  <progress value={60} max={100}></progress> Wants
             </div>
          </div>
      </div>
      <div className={styles.row}>
          <div className={styles.progressBarContainer}>
           Needs
            <div>
                  <progress value={100} max={100}></progress> Rent
            </div>
            <div>
                  <progress value={100} max={100}></progress> Groceries
            </div>
              <div>
                  <progress value={100} max={100}></progress> Car Payment
             </div>
              <div>
                  <progress value={100} max={100}></progress> Insurance
             </div>
          </div>
          <div className={styles.progressBarContainer}>
            Debt Repayment
            <div>
                  <progress value={100} max={100}></progress> Student Loan 
            </div>
            <div>
                  <progress value={100} max={100}></progress> Credit Cards
            </div>
              <div>
                  <progress value={100} max={100}></progress> Personal Loan
             </div>
          </div>
           <div className={styles.progressBarContainer}>
            Wants
            <div>
                  <progress value={45} max={100}></progress> Restaurants
            </div>
            <div>
                  <progress value={100} max={100}></progress> Credit Cards
            </div>
              <div>
                  <progress value={100} max={100}></progress> Personal Loan
             </div>
          </div>
      </div>
    </div>
  )

}

export default Dashboard;
