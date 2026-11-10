import { animate, stagger } from "../motion.js";

export default function (el) {
  const [middle, left, right] = el.querySelectorAll("g");
  animate([left, middle, right], { rotate: [0, -12, 6, 0] }, { duration: 0.3, delay: stagger(0.09) });
}
