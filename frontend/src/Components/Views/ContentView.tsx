import React, { useState } from 'react';
import BudgetSetup from '../Budget/BudgetSetup';
import styles from './views.module.scss';
import Dashboard from '../Dashboard/Dashboard';


interface ContentViewProps {
  activeView: string;
}



const ContentView = ({ activeView }: ContentViewProps) => {

    return (
        <div className={`${styles['col-md-11']}`}>
            {activeView === 'budget' && <BudgetSetup />}
            {activeView === 'dashboard' && <Dashboard />}
        </div>
    )

}
export default ContentView;