import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cabinet = el.querySelector("rect");
  animate(cabinet, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const handles = ["M10 6h4", "M10 16h4"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(handles, { scaleX: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
