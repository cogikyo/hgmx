import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const pad = el.querySelector("path");
  animate(pad, { scale: [1, 0.85, 1.1, 1] }, { duration: τmid });

  const toes = el.querySelectorAll("circle");
  animate(toes, { scale: [0, 1.5, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
