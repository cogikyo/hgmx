import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = ["M9 8h7", "M8 12h6", "M11 16h5"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
