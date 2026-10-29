import { animate, stagger } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll("path");
  animate(lines, { scaleX: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
