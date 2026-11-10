import { animate, stagger } from "../motion.js";

export default function (el) {
  const steps = el.querySelectorAll("rect");
  animate(steps, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
