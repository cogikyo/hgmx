import { animate, stagger } from "../motion.js";

export default function (el) {
  const arms = el.querySelectorAll("g");
  animate(arms, { scale: [1, 1.25, 0.95, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
