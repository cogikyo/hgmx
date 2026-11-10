import { animate, stagger } from "../motion.js";

export default function (el) {
  const steps = el.querySelectorAll('path[d="M3 17h18"], path[d="M8 17v-4h13"], path[d="M20.77 9H12v4"]');
  animate(steps, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
