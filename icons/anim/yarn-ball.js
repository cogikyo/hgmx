import { animate, stagger } from "../motion.js";

export default function (el) {
  const strands = el.querySelectorAll("path");
  animate(strands, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
