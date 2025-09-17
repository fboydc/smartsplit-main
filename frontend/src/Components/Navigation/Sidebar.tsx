import React, { useContext } from "react";
import styles from "./navigation.module.scss";

// Import your PNG icons
import budgetIcon from "../../assets/icons/budget.png";
import dashboardIcon from "../../assets/icons/dashboard.png";
import bankIcon from "../../assets/icons/bank.png";
import transactionsIcon from "../../assets/icons/transactions.png";


interface SidebarProps {
  activeView: string;
  setActiveView: React.Dispatch<React.SetStateAction<string>>;
}



const Sidebar = ({ activeView, setActiveView }: SidebarProps) => {
    return (
            <div className={`${styles.sidebar} ${styles['sidebar-nav']}`}>
                <ul className={`${styles.listGroup} ${styles.siderbarNav}`}>
                    <li className={styles.listGroupItem}>
                            <button className={`${styles['btn']} ${ activeView === 'budget' ? styles['active'] : ''}`} onClick={() => setActiveView('budget')}>
                                <img src={budgetIcon} alt="Budget Distribution" width={50} height={50} />
                            </button>
                    </li>
                    <li className={styles.listGroupItem}> 
                            <button className={`${styles['btn']} ${activeView === 'dashboard' ? styles['active']: ""}`}  onClick={() => setActiveView('dashboard')}>
                                 <img src={dashboardIcon} alt="Dashboard" width={50} height={50} />
                            </button>
                    </li>
                    <li className={styles.listGroupItem}>
                        <button className={`${styles['btn']}`}>
                             <img src={bankIcon} alt="Bank Setup" width={50} height={50} />
                        </button>
                    </li>
                    <li className={styles.listGroupItem}> 
                        <button className={`${styles['btn']}`}>
                             <img src={transactionsIcon} alt="Transactions" width={50} height={50} />
                        </button>
                    </li>
                </ul>
            </div>
    )

}

export default Sidebar;