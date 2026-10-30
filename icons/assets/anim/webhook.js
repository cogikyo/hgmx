import { animate, stagger } from "../motion.js";

export default function (el) {
  const hooks = el.querySelectorAll("path");
  animate(hooks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
