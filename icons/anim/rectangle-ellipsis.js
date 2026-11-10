import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const rect = el.querySelector("rect");
  animate(rect, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dots = ["M7", "M12", "M17"].map((x) => el.querySelector(`path[d="${x} 12h.01"]`));
  animate(dots, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
