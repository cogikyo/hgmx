import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  const code = [el.querySelector("rect"), ...["M17 8V7", "M17 12v4", "M7 17h.01"].map((d) => el.querySelector(`path[d^="${d}"]`))];
  animate(code, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
