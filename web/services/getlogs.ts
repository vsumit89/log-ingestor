import axios from "axios";
import { SearchFilters } from "@/commons/types/searchFilter";

export const getLogs = async (searchQuery: string, filters: SearchFilters) => {
  const response = await axios.get(`http://localhost:3000`, {
    params: {
      q: searchQuery,
      ...filters,
    },
  });
  return response.data;
};
