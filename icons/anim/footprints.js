import { animate, stagger } from "../motion.js";

export default function (el) {
  const feet = el.querySelectorAll("g");
  animate(feet, { translateY: [0, -2, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
