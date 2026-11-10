import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const core = el.querySelector("circle");
  animate(core, { scale: [1, 0.85, 1] }, { duration: τmid });

  const dots = el.querySelectorAll("path");
  animate(dots, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
