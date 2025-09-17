import React, { useEffect, useState, ChangeEvent } from 'react';
import styles from "./ExpandableTable.module.scss";
import { format } from 'path';
import {AllocationGroup, Expense } from "../../models/types";
import { v4 as uuidv4} from 'uuid';


interface ExpandableTableProps {
  fields: Expense[]; // Array of Expense objects
  setFields: (updatedFields: Expense[]) => void; // Function to update the fields
  categories: { id: string; name: string }[]; // List of categories for the dropdown
  formatCurrency: (value: string) => string; // Function to format currency
  subtotal: string; // Subtotal to display 
}


const ExpandableTable: React.FC<ExpandableTableProps> =({
    fields,
    setFields,
    categories,
    formatCurrency,
    subtotal,
}) => {
   

    const addRow = () => {

        //var row: Expense = {id: fields.length + 1, name: name, amount: amount, category: categoryId}
        
          const newRow: Expense = {
                id: uuidv4(),
                description: "",
                amount: 0,
                category: "",
                allocation_type: "", // Default value for allocation_type
           };

        setFields([...fields, newRow]);
    }

    const removeRow = (id: string) => {
        const newRows = fields.filter((row) => row.id !== id);
        setFields(newRows);
    }

    const isValidNumber = (value: string) => {
        const numberRegex = /^\d*\.?\d*$/
        return numberRegex.test(value);
    }
    
    const handleChange = (id: string, field: keyof Expense, value: string) => {

        if (field === "amount" && !isValidNumber(value)) {
            return
        }

        const updatedFields = fields.map((row) =>
        row.id === id
            ? {
                ...row,
                [field]: value,
            }
            : row
    );
    setFields(updatedFields);

    }
   
    return (
        <div>
            <table className={styles.table}>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Amount</th>
                        <th>Category</th>
                        <th></th>
                    </tr>
                </thead>
                
                <tbody>
                    {
                        fields && fields.map((row) => (
                            <tr key={row.id}>
                                <td>
                                    <input type="text" value={row.description} onChange={(e)=> handleChange(row.id, "description", e.target.value)}/>
                                </td>
                                <td>
                                     <input type="text" value={row.amount} onChange={(e)=> handleChange(row.id, "amount", e.target.value)}/>
                                </td>
                                <td>
                                <select value={row.category || "" } onChange={(e)=> handleChange(row.id, "category", e.target.value)}>
                                    {categories.map((category) => (
                                        <option key={category.id} value={category.id}>
                                            {category.name}
                                        </option>
                                    ))}
                                </select>
                                </td>
                                <td>
                                    <button className={styles.tableButton} onClick={() => removeRow(row.id)}>Remove</button>
                                </td>
                            </tr>
                        ))

                    /*}
                    {console.log("fields: " + JSON.stringify(fields))}
                    {fields && fields.map((row) => (
                        <tr key={row.id}>
                            <td>
                                <input type="text" value={row.description} onChange={(e)=> handleChange(row.id, "name", e.target.value)}/>
                            </td>
                            <td>
                                 <input type="text" value={row.amount} onChange={(e)=> handleChange(row.id, "amount", e.target.value)}/>
                            </td>
                            <td>

                                <select value={row.category || "" } onChange={(e)=> handleChange(row.id, "category", e.target.value)}>
                                {categories.map((category) => (
                                    <option key={category.id} value={category.id}>
                                        {category.name}
                                    </option>
                                ))}
                                </select>
                            </td>
                            <td>
                                <button className={styles.tableButton} onClick={() => removeRow(row.id)}>Remove</button>
                            </td>
                        </tr>
                    ))}
                    <tr className={styles.subTotalRow}>
                        <td>
                            <p>Sub-Total</p>
                        </td>
                        <td colSpan={6}>
                            <p>{subtotal}</p>
                        </td>
                    </tr>
                    <tr>
                        <td colSpan={3}>
                            <button className={styles.tableButton} onClick={() => addRow({name: "", amount: "", categoryId: ""})}>Add Row</button>
                        </td>
                        <td />
                    </tr>
                    */}
                    
                       <tr className={styles.subTotalRow}>
                            <td>
                                <p>Sub-Total</p>
                            </td>
                            <td colSpan={6}>
                                <p>{subtotal}</p>
                            </td>
                        </tr>
                         <tr>
                            <td colSpan={3}>
                                <button className={styles.tableButton} onClick={() => addRow()}>Add Row</button>
                            </td>
                            <td />
                        </tr>
                </tbody>
            </table>
        </div>
    )


}


export default ExpandableTable;