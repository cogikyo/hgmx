import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const mic = el.querySelector("rect");
  animate(mic, { scale: [1, 0.85, 1] }, { duration: τmid });

  const arcs = ["M6 11", "M18 11", "M2 11", "M22 11"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(arcs, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
