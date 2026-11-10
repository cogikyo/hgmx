import { animate, stagger } from "../motion.js";

export default function (el) {
  const digits = ["m10 4", "M8 22", "M16 14"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(digits, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
