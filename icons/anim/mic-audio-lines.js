import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const mic = el.querySelector("rect");
  animate(mic, { scale: [1, 0.85, 1] }, { duration: τmid });

  const bars = ["M6 6", "M18 5", "M2 10", "M22 10"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(bars, { scaleY: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
