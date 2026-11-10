import { animate, stagger } from "../motion.js";

export default function (el) {
  const rows = el.querySelectorAll('path[d="m2 22 5-8"], path[d="m7 22 5-8"], path[d="M22 22H12l5-8"], path[d="M15 18h7"]');
  animate(rows, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
