export interface Expense {
    id: string;
    description: string;
    amount: number;
    category: string;
    allocation_type: string;
  }
  
  export interface Income {
    id: string;
    description: string;
    amount: number;
    frequency: string;
  }
  
  export interface Allocation {
    type: string;
    description: string;
    factor: number;
  }

  export interface AllocationGroup {
    allocation_type: string;
    allocation_total: number;
    current_allocation: number;
    allocation_pct: number;
    expenses: Expense[];
  }