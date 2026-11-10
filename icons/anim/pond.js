import { animate, stagger } from "../motion.js";

export default function (el) {
  const cattails = el.querySelectorAll("g");
  animate(cattails, { rotate: [0, -10, 6, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
