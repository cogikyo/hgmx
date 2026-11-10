import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const speakers = el.querySelectorAll("circle");
  animate(speakers, { scale: [1, 1.5, 0.85, 1.25, 1] }, { duration: τmax, delay: stagger(0.06) });
}
