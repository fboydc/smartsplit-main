import React, { useEffect, useState, ChangeEvent } from 'react';
import styles from "./ExpandableTable.module.scss";
import { format } from 'path';

interface StaticTableProps {
    fields: {
        headings: string[],
        rows: {
          columns: string[];
        }[];
      };
}

const StaticTable: React.FC<StaticTableProps> = ({ fields }) => {
   
   
    return (
        <div>
            <table className={styles.table}>
                <thead>
                    <tr>
                        {fields.headings.length > 0 && fields.headings.map((heading, index) => (
                            <th key={index}>{heading}</th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                {fields.rows.map((row, rowIndex) => (
                        <tr key={rowIndex}>
                        {row.columns.map((column, colIndex) => (
                            <td key={colIndex} colSpan={colIndex}>{column}</td>
                        ))}
                        </tr>
                ))}
                </tbody>
            </table>
        </div>
    )


}


export default StaticTable;