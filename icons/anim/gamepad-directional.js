import { animate, stagger } from "../motion.js";

export default function (el) {
  const [down, right, left, up] = el.querySelectorAll("path");
  animate([up, right, down, left], { scale: [1, 0.75, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
