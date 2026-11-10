import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const columns = ["M8 7v7", "M12 7v4", "M16 7v9"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(columns, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
