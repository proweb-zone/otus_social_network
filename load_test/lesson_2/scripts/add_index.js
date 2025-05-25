import http from 'k6/http';
import { check, group } from 'k6';

export let options = {
   stages: [
       { duration: '3m', target: 1 },
     ],
};

export default function () {
   group('API uptime check', () => {
       const response = http.get('http://localhost:3002/user/search/лев абра');
       check(response, {
           "status code should be 200": res => res.status === 200,
       });
   });
};
