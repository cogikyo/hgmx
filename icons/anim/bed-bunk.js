import { animate, stagger } from "../motion.js";

export default function (el) {
  const bunks = el.querySelectorAll("g");
  animate(bunks, { scaleY: [1, 0.75, 1.25, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
