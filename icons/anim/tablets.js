import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const scores = el.querySelectorAll("path");
  animate(scores, { rotate: [0, 180] }, { duration: τmid, delay: stagger(0.12) });
}
