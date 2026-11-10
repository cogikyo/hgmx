import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const panel = el.querySelector("rect");
  animate(panel, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const screws = ["M7 7", "M17 7", "M17 17", "M7 17"].map((d) => el.querySelector(`path[d="${d}h.01"]`));
  animate(screws, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
