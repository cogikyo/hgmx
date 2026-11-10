import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const button = el.querySelector("circle");
  animate(button, { scale: [1, 0.7, 1] }, { duration: τmid });

  const dots = ["M8", "M12", "M16"].map((x) => el.querySelector(`path[d="${x} 8h.01"]`));
  animate(dots, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
